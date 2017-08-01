package mediationcontainer

import(
	"testing"
	goproto "github.com/golang/protobuf/proto"
)

type myMsg struct {
	name string
}

func (m *myMsg) Reset() { *m = myMsg{}}
func (m *myMsg) String() string { return m.name }
func (m *myMsg) ProtoMessage() {}

type testHandler struct {
	id int
}

func (h *testHandler)HandleRawMessage(msg []byte) (goproto.Message, error) {
	return &myMsg{name: "hello"}, nil
}

func TestPipeline_Push(t *testing.T) {
	pipe := NewPipeline()

	h := &testHandler{ id: 1}
	pipe.Push(h)

	if pipe.Len() != 1 {
		t.Error("Push test failed.")
	}
}

func TestPipeline_Pop(t *testing.T) {
	pipe := NewPipeline()

	h := &testHandler{ id: 1}
	pipe.Push(h)

	n, err := pipe.Pop()
	if err != nil {
		t.Error("Pop test failed.")
	}

	h2, ok := n.(*testHandler)
	if !ok {
		t.Error("Pop test failed.")
	}

	if h2.id != h.id {
		t.Error("Pop test failed.")
	}
}


func TestPipeline_Peek(t *testing.T) {
	pipe := NewPipeline()

	h := &testHandler{ id: 1}
	pipe.Push(h)

	n, err := pipe.Peek()
	if err != nil {
		t.Error("Peek test failed.")
	}

	if pipe.Len() != 1 {
		t.Error("Peek test failed.")
	}

	h2, ok := n.(*testHandler)
	if !ok {
		t.Error("Peek test failed.")
	}

	if h2.id != h.id {
		t.Error("Peek test failed.")
	}
}

func TestPipeline_Push2(t *testing.T) {
	pipe := NewPipeline()

	n := 10
	for i := 0; i < n; i ++ {
		h := &testHandler{ id: i}
		pipe.Push(h)
	}

	if pipe.Len() != n {
		t.Error("Push many test failed.")
	}
}


func TestPipeline_Pop2(t *testing.T) {
	pipe := NewPipeline()

	n := 10
	for i := 0; i < n; i ++ {
		h := &testHandler{ id: i}
		pipe.Push(h)
	}

	for i := 0; i < n; i ++ {
		h, err := pipe.Pop()
		if err != nil {
			t.Error("Pop many test failed: %v", err)
		}

		hh, ok := h.(*testHandler)
		if !ok {
			t.Error("Pop many test failed.")
		}

		if hh.id != i {
			t.Errorf("Pop many test failed [%v Vs. %v].", hh.id, i)
		}
	}

	if pipe.Len() != 0 {
		t.Error("Pop many test failed.")
	}
}
