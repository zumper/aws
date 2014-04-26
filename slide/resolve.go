package aws

func Resolve(wsdl WSDL) (req, resp map[string][]string) {
	req, resp = make(map[string][]string), make(map[string][]string)
	for _, op := range wsdl.Operation {
		inName := Unqualify(op.Input.Message)
		msgIn := wsdl.Message[inName]
		inName = ResolveMsg(inName, wsdl)
		req[inName] = []string{} // Some types have no deps
		for _, r := range resolve(Unqualify(msgIn.Part.Element), wsdl) {
			req[inName] = append(req[inName], r)
		}

		outName := Unqualify(op.Output.Message)
		msgOut := wsdl.Message[outName]
		outName = ResolveMsg(outName, wsdl)
		resp[outName] = []string{} // Some types have not deps
		for _, r := range resolve(Unqualify(msgOut.Part.Element), wsdl) {
			resp[outName] = append(resp[outName], r)
		}
	}
	return
}
