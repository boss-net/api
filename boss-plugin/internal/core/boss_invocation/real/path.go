package real

func (r *RealBackwardsInvocation) bossPath(path ...string) string {
	path = append([]string{"inner", "api"}, path...)
	return r.bossInnerApiBaseurl.JoinPath(path...).String()
}
