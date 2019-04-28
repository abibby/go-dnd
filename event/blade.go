package event

import (
	"fmt"

	"github.com/zwzn/dnd/blade"
)

func (r *recharge) String() string {
	eventType := ""
	switch r.event {
	case LongRest:
		eventType = "long rest"
	case ShortRest:
		eventType = "short rest"
	default:
		eventType = string(r.event)
	}

	used := ""
	if r.use == r.total {
		if r.total == r.current {
			used = "available"
		} else {
			used = "used"
		}
	} else {
		used = fmt.Sprintf("%d of %d available", r.current/r.use, r.total/r.use)
	}
	if r.use != 1 {
		return fmt.Sprintf("Recharges after %d %ss, %s", r.use, eventType, used)
	}
	return fmt.Sprintf("Recharges after a %s, %s", eventType, used)
}
func (c *chWrapper) updatePreBlade() {
	b := blade.New()
	b.Directive("recharge", func(args blade.Args) string {
		var name string
		var event EventType
		count := 1
		use := 1
		err := args.Unmarshal(&name, &event, &count, &use)
		if err != nil {
			return err.Error()
		}

		count *= use
		c.recharge[name] = recharge{
			event:   event,
			total:   count,
			use:     use,
			current: count,
		}
		return fmt.Sprintf("@recharge(%s)", args)
	})
	c.Blade(b)
}

func (c *chWrapper) updatePostBlade() {
	b := blade.New()
	b.Directive("recharge", func(args blade.Args) string {
		var name string
		err := args.Unmarshal(&name)
		if err != nil {
			return err.Error()
		}
		r := c.recharge[name]
		return r.String()
	})
	b.Directive("mod", func(args blade.Args) string {
		var mod string
		err := args.Unmarshal(&mod)
		if err != nil {
			return err.Error()
		}

		return sign(c.AbilityScoreMods()[mod])
	})
	b.Directive("save", func(args blade.Args) string {
		var mod string
		err := args.Unmarshal(&mod)
		if err != nil {
			return err.Error()
		}

		return sign(c.SavingThrows()[mod])
	})
	c.Blade(b)
}

func sign(i int) string {
	if i >= 0 {
		return fmt.Sprintf("+%d", i)
	}
	return fmt.Sprint(i)
}
