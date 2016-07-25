package controller

import "gosaic/environment"

func PartialAspect(env environment.Environment, columns, rows int) error {
	aspectService, err := env.AspectService()
	if err != nil {
		return err
	}

	aspect, err := aspectService.FindOrCreate(columns, rows)
	if err != nil {
		return err
	}

	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		return err
	}

	for {
		gidxs, err := gidxPartialService.FindMissing(aspect, "gidx.id ASC", 1000, 0)
		if err != nil {
			return err
		}
		if len(gidxs) == 0 {
			break
		}

		for _, gidx := range gidxs {
			env.Printf("Creating partial for index %d\n", gidx.Id)
			_, err := gidxPartialService.Create(gidx, aspect)
			if err != nil {
				return err
			}
		}
	}

	return nil
}