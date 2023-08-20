package requests

type ReqCB func(req *SRequest) error
type RespCB func(req *SResponse) error

type reqMid struct {
	curCB ReqCB
	next  *reqMid
}

func NewReqMid(f ReqCB) *reqMid {
	return &reqMid{
		curCB: f,
	}
}

func (r *reqMid) add(f ReqCB) {
	if r.next != nil {
		r.next.add(f)
	}
	r.next = NewReqMid(f)
}

func (r *reqMid) run(req *SRequest) error {
	err := r.curCB(req)
	if err != nil {
		return err
	}
	if r.next != nil {
		return r.next.run(req)
	}
	return nil
}

type respMid struct {
	curCB RespCB
	next  *respMid
}

func NewRespMid(f RespCB) *respMid {
	return &respMid{
		curCB: f,
	}
}

func (r *respMid) add(f RespCB) {
	if r.next != nil {
		r.next.add(f)
	}
	r.next = NewRespMid(f)
}

func (r *respMid) run(resp *SResponse) error {
	err := r.curCB(resp)
	if err != nil {
		return err
	}
	if r.next != nil {
		return r.next.run(resp)
	}
	return nil
}
