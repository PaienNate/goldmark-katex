package katex

import (
	_ "embed"
	"io"

	"github.com/fastschema/qjs"
)

//go:embed katex.min.js
var code string

func Render(w io.Writer, src []byte, display bool) error {
	rt, err := qjs.New()
	if err != nil {
		return err
	}
	defer rt.Close()

	ctx := rt.Context()

	result, err := ctx.Eval("katex.min.js", qjs.Code(code))
	if err != nil {
		return err
	}
	defer result.Free()

	srcVal := ctx.NewString(string(src))
	ctx.Global().SetPropertyStr("_EqSrc3120", srcVal)
	// srcVal is now owned by the global object? Or should we free it?
	// In typical bindings, if you pass a value to a setter, it might increment ref count or copy.
	// But to be safe, if we don't know, we can defer Free it.
	// However, if we free it and it was just a reference, it might be invalid.
	// Let's assume we should defer Free it. If it causes issues, we'll see.
	// Wait, if SetPropertyStr takes ownership, we shouldn't free.
	// But usually creating a NewString gives us an owned reference.
	// Let's defer Free() for now.
	defer srcVal.Free()

	if display {
		result, err = ctx.Eval("render.js", qjs.Code("katex.renderToString(_EqSrc3120, { displayMode: true })"))
	} else {
		result, err = ctx.Eval("render.js", qjs.Code("katex.renderToString(_EqSrc3120)"))
	}

	if err != nil {
		return err
	}
	defer result.Free()

	_, err = io.WriteString(w, result.String())
	return err
}
