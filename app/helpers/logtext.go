package helpers

import "regexp"
import "fmt"

// Similar to ruby gsub!
func gsub(text *string, pattern string, repl string) (error) {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return err
    }
    *text = string(re.ReplaceAll([]byte(*text), []byte(repl)))
    return nil
}

func ProcessLogText(log string) (string, error) {

    // Linies llargues
    err := gsub(&log, "\\S{45}", "$1\n")
    if err != nil {
        return "", err
    }

    // Nicks i hores
    principi := "(^|\\n)";
    hora := "[\\[\\(] *([0-9]{2}\\:[0-9]{2}\\:?[0-9]{0,2})[ap]?m? *[\\]\\)]";
    nick := "(<|[\\[\\( ])? *[@+]?([^0-9 .:@]{1}[^ :@]*) *(\\): |>|: |[\\]\\)])";

    // Hora + nick
    var rex string
    rex = fmt.Sprintf("(?m)%v%v *%v *", principi, hora, nick)
	err = gsub(&log, rex, "$1<span class=\"hora\">[$2]</span> <span class=\"nick_deco\">&lt;</span><span class=\"nick\">$4</span><span class=\"nick_deco\">&gt;</span> ")
    if err != nil {
        return "", err
    }

	// Hora sola
    rex = fmt.Sprintf("(?m)%v%v *", principi, hora)
	gsub(&log, rex, "$1<span class=\"hora\">[$2]</span> ")

	// Nick sol
    rex = fmt.Sprintf("(?m)%v%v *", principi, nick)
	gsub(&log, rex, "$1<span class=\"nick_deco\">&lt;</span><span class=\"nick\">$3</span><span class=\"nick_deco\">&gt;</span> ")

    // NL 2 BR
    gsub(&log, "\\n", "<br>")

    // TODO: Highlight

    // # Highlight
    // unless paraules.nil?
    //   for paraula in paraules
    //     log.gsub! /([^[:alnum:]]|^)(#{paraula})([^[:alnum:]]|$)/i, "\\1<span class=\"highlight\">\\2</span>\\3"
    //   end
    // end
    //
    return log, nil
}
