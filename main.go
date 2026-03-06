package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calc1", calc1Page)
	http.HandleFunc("/calc1/result", calc1Result)

	http.HandleFunc("/calc2", calc2Page)
	http.HandleFunc("/calc2/result", calc2Result)

	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func calc1Page(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "calc1.html", nil)
}

func calc2Page(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "calc2.html", nil)
}

func round(val float64, decimals int) float64 {
	format := "%." + strconv.Itoa(decimals) + "f"
	v, _ := strconv.ParseFloat(fmt.Sprintf(format, val), 64)
	return v
}

func calc1Result(w http.ResponseWriter, r *http.Request) {

	hp, _ := strconv.ParseFloat(r.FormValue("hp"), 64)
	cp, _ := strconv.ParseFloat(r.FormValue("cp"), 64)
	sp, _ := strconv.ParseFloat(r.FormValue("sp"), 64)
	np, _ := strconv.ParseFloat(r.FormValue("np"), 64)
	op, _ := strconv.ParseFloat(r.FormValue("op"), 64)
	wp, _ := strconv.ParseFloat(r.FormValue("wp"), 64)
	ap, _ := strconv.ParseFloat(r.FormValue("ap"), 64)

	krs := 100 / (100 - wp)
	krg := 100 / (100 - wp - ap)

	hS := hp * krs
	cS := cp * krs
	sS := sp * krs
	nS := np * krs
	oS := op * krs
	aS := ap * krs

	hG := hp * krg
	cG := cp * krg
	sG := sp * krg
	nG := np * krg
	oG := op * krg

	qrn := (339*cp + 1030*hp - 108.8*(op-sp) - 25*wp) / 1000
	qrg := (qrn + 0.025*wp) * (100 / (100 - wp - ap))
	qd := (qrn + 0.025*wp) * (100 / (100 - wp))

	result := fmt.Sprintf(`
Коефіцієнт переходу від робочої до сухої маси = %.3f
Коефіцієнт переходу від робочої до горючої маси = %.3f

Склад сухої маси палива:
Hc = %.2f %%
Cc = %.2f %%
Ss = %.2f %%
Ns = %.2f %%
Os = %.2f %%
As = %.2f %%

Склад горючої маси палива:
Hc = %.2f %%
Cc = %.2f %%
Ss = %.2f %%
Ns = %.2f %%
Os = %.2f %%

Нижча теплота згоряння:
Робоча маса = %.3f МДж/кг
Суха маса = %.3f МДж/кг
Горюча маса = %.3f МДж/кг
`,
		round(krs, 3), round(krg, 3),
		round(hS, 2), round(cS, 2), round(sS, 2),
		round(nS, 2), round(oS, 2), round(aS, 2),
		round(hG, 2), round(cG, 2), round(sG, 2),
		round(nG, 2), round(oG, 2),
		round(qrn, 3), round(qd, 3), round(qrg, 3))

	templates.ExecuteTemplate(w, "result.html", result)
}

func calc2Result(w http.ResponseWriter, r *http.Request) {

	hg, _ := strconv.ParseFloat(r.FormValue("hg"), 64)
	cg, _ := strconv.ParseFloat(r.FormValue("cg"), 64)
	sg, _ := strconv.ParseFloat(r.FormValue("sg"), 64)
	og, _ := strconv.ParseFloat(r.FormValue("og"), 64)
	vg, _ := strconv.ParseFloat(r.FormValue("vg"), 64)
	wg, _ := strconv.ParseFloat(r.FormValue("wg"), 64)
	ag, _ := strconv.ParseFloat(r.FormValue("ag"), 64)
	x, _ := strconv.ParseFloat(r.FormValue("x"), 64)

	solut1 := (100 - wg - ag) / 100
	solut2 := (100 - wg) / 100

	cg2 := cg * solut1
	hg2 := hg * solut1
	og2 := og * solut1
	sg2 := sg * solut1
	ag2 := ag * solut2
	vg2 := vg * solut2

	q := x*((100-wg-ag2)/100) - 0.025*ag

	result := fmt.Sprintf(`
Склад робочої маси мазуту становитиме:

Hp = %.2f %%
Cp = %.2f %%
Sp = %.2f %%
Op = %.2f %%
Vp = %.2f %%
Ap = %.2f %%

Нижча теплота згоряння мазуту та робочої маси за заданим складом компонентів палива становить:

%.2f %%
`,
		round(hg2, 2), round(cg2, 2), round(sg2, 2),
		round(og2, 2), round(vg2, 2), round(ag2, 2),
		round(q, 2))

	templates.ExecuteTemplate(w, "result.html", result)
}
