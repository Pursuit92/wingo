package frame

import (
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/BurntSushi/wingo/logger"
)

type Parent struct {
	*xwindow.Window
	isMapped    bool
	MoveState   *MoveState
	ResizeState *ResizeState
}

type MoveState struct {
	Moving       bool
	RootX, RootY int
}

type ResizeState struct {
	Resizing            bool
	RootX, RootY        int
	X, Y, Width, Height int
	Xs, Ys, Ws, Hs      bool
}

func newParent(X *xgbutil.XUtil, cid xproto.Window) (*Parent, error) {
	parent, err := xwindow.Generate(X)
	if err != nil {
		logger.Error.Printf("Could not create a parent window for client "+
			"with id '%d' because: %s", cid, err)
		logger.Error.Fatalf("In a state where no new windows can be created. " +
			"Unfortunately, we must exit.")
	}

	err = parent.CreateChecked(X.RootWin(), 0, 0, 1, 1,
		xproto.CwEventMask,
		xproto.EventMaskSubstructureRedirect|
			xproto.EventMaskButtonPress|
			xproto.EventMaskButtonRelease|
			xproto.EventMaskFocusChange)
	if err != nil {
		return nil, err
	}

	err = xproto.ReparentWindowChecked(X.Conn(),
		cid, parent.Id, 0, 0).Check()
	if err != nil {
		return nil, err
	}

	return &Parent{
		Window:      parent,
		MoveState:   &MoveState{},
		ResizeState: &ResizeState{},
		isMapped:    false,
	}, nil
}

func (par *Parent) Deparent(cid xproto.Window) {
	err := xproto.ReparentWindowChecked(
		par.X.Conn(), cid, par.X.RootWin(), 0, 0).Check()
	if err != nil {
		logger.Warning.Printf("Could not deparent client window: %s", err)
	}
}

func (par *Parent) Map() {
	par.Window.Map()
	par.isMapped = true
}

func (par *Parent) Unmap() {
	par.Window.Unmap()
	par.isMapped = false
}
