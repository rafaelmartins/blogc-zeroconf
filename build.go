package main

func build(ctx *context, out string) error {
	bctxs, err := ctx.getBuildCtxs(out, true)
	if err != nil {
		return err
	}

	for _, c := range bctxs {
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			return err
		}
	}

	return nil
}
