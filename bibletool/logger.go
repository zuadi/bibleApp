package bibletool

func (bt *Bibletool) LogInfo(caller, msg string) {
	bt.Logger.Info(caller, msg)
}

func (bt *Bibletool) LogWarning(caller, msg string) {
	bt.Logger.Warning(caller, msg)
}

func (bt *Bibletool) LogError(caller string, msg any) {
	bt.Logger.Error(caller, msg)
}
