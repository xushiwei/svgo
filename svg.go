// svg: generate SVG objects

package svg

import (
	"fmt"
	"os"
	"xml"
	"strings"
)

const svginit = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink"
     width="%d" height="%d">
`


// Structure and Metadata

func Start(w int, h int)  { fmt.Printf(svginit, w, h) }
func End()                { fmt.Println("</svg>") }
func Gstyle(s string)     { fmt.Println(group("style", s)) }
func Gtransform(s string) { fmt.Println(group("transform", s)) }
func Gid(s string)        { fmt.Println(group("id", s)) }
func Gend()               { fmt.Println("</g>") }
func Def()                { fmt.Println("<defs>") }
func DefEnd()             { fmt.Println("</defs>") }
func Desc(s string)       { tt("desc", "", s) }
func Title(s string)      { tt("title", "", s) }
func Use(x int, y int, link string, s ...string) {
	fmt.Printf(`<use %s %s %s`, loc(x, y), href(link), endstyle(s))
}

// Shapes

func Circle(x int, y int, r int, s ...string) {
	fmt.Printf(`<circle cx="%d" cy="%d" r="%d" %s`, x, y, r, endstyle(s))
}

func Ellipse(x int, y int, w int, h int, s ...string) {
	fmt.Printf(`<ellipse cx="%d" cy="%d" rx="%d" ry="%d" %s`,
		x, y, w, h, endstyle(s))
}

func Polygon(x []int, y []int, s ...string) { poly(x, y, "polygon", s) }

func Rect(x int, y int, w int, h int, s ...string) {
	fmt.Printf(`<rect %s %s`, dim(x, y, w, h), endstyle(s))
}

func Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string) {
	fmt.Printf(`<rect %s rx="%d" ry="%d" %s`, dim(x, y, w, h), rx, ry, endstyle(s))
}

func Square(x int, y int, s int, style ...string) {
	Rect(x, y, s, s, style)
}

// Curves

func Arc(sx int, sy int, ax int, ay int, r int, large bool, sweep bool, ex int, ey int, s ...string) {
	fmt.Printf(`%s A%s %d %s %s %s" %s`,
		ptag(sx, sy), coord(ax, ay), r, onezero(large), onezero(sweep), coord(ex, ey), endstyle(s))
}

func Bezier(sx int, sy int, cx int, cy int, px int, py int, ex int, ey int, s ...string) {
	fmt.Printf(`%s C%s %s %s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(px, py), coord(ex, ey), endstyle(s))
}

func Qbezier(sx int, sy int, cx int, cy int, ex int, ey int, tx int, ty int, s ...string) {
	fmt.Printf(`%s Q%s %s T%s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(ex, ey), coord(tx, ty), endstyle(s))
}

// Lines

func Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	fmt.Printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" %s`, x1, y1, x2, y2, endstyle(s))
}

func Polyline(x []int, y []int, s ...string) { poly(x, y, "polyline", s) }


// Image

func Image(x int, y int, w int, h int, link string, s ...string) {
	fmt.Printf("<image %s %s %s", dim(x, y, w, h), href(link), endstyle(s))
}

// Text

func Text(x int, y int, t string, s ...string) {
	if len(s) > 0 {
		tt("text", " "+loc(x, y)+" "+style(s[0]), t)
	} else {
		tt("text", " "+loc(x, y)+" ", t)
	}
}

// Color

func RGB(r int, g int, b int) string { return fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b) }
func RGBA(r int, g int, b int, a float) string {
	return fmt.Sprintf(`fill-opacity:%.2f; %s`, a, RGB(r, g, b))
}

// Utility

func Grid(x int, y int, w int, h int, n int, s ...string) {

	if len(s) > 0 {
		Gstyle(s[0])
	}
	for ix := x; ix <= x+w; ix += n {
		Line(ix, y, ix, y+h)
	}

	for iy := y; iy <= y+h; iy += n {
		Line(x, iy, x+w, iy)
	}
	if len(s) > 0 {
		Gend()
	}

}

// Support functions

func style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}

func pp(x []int, y []int, tag string) {
	if len(x) != len(y) {
		return
	}
	fmt.Print(tag)
	for i := 0; i < len(x); i++ {
		fmt.Print(coord(x[i], y[i]) + " ")
	}
}

func endstyle(s []string) string {
	if len(s) > 0 {
		if strings.Index(s[0], "=") > 0 {
			return s[0] + "/>\n"
		} else {
			return style(s[0]) + "/>\n"
		}
	}
	return "/>\n"
}

func tt(tag string, attr string, s string) {
	fmt.Print("<" + tag + attr + ">")
	xml.Escape(os.Stdout, []byte(s))
	fmt.Println("</" + tag + ">")
}

func poly(x []int, y []int, tag string, s ...string) {
	pp(x, y, "<"+tag+` points="`)
	fmt.Print(`" ` + endstyle(s))
}

func onezero(flag bool) string {
	if flag {
		return "1"
	}
	return "0"
}

func coord(x int, y int) string { return fmt.Sprintf(`%d,%d`, x, y) }
func ptag(x int, y int) string  { return fmt.Sprintf(`<path d="M%s`, coord(x, y)) }
func loc(x int, y int) string   { return fmt.Sprintf(`x="%d" y="%d"`, x, y) }
func href(s string) string      { return fmt.Sprintf(`xlink:href="%s"`, s) }
func dim(x int, y int, w int, h int) string {
	return fmt.Sprintf(`x="%d" y="%d" width="%d" height="%d"`, x, y, w, h)
}
func group(tag string, value string) string { return fmt.Sprintf(`<g %s="%s">`, tag, value) }
