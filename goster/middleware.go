package goster

func HandleLog(c *Ctx, g *Goster, err error) {
	m := c.Request.Method
	u := c.Request.URL.String()

	if err != nil {
		l := err.Error()
		g.Logs = append(g.Logs, l)
		LogError(l, g.Logger)
		return
	}
	l := "[" + m + "]" + " ON ROUTE " + u
	g.Logs = append(g.Logs, l)
	LogInfo(l, g.Logger)
}
