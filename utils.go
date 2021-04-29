package animeworld

import (
	"errors"
	"github.com/KiritoNya/htmlutils"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

func equalSliceFloat(a, b []float64) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func getServers(node *html.Node) ([]Server, error) {

	var s []Server

	servers, err := htmlutils.QuerySelector(node, "span", "class", "tab server-tab")
	if err != nil {
		return nil, errors.New("Error to get servers: " + err.Error())
	}

	for _, server := range servers {

		var srv Server

		id, err := htmlutils.GetValueAttr(server, "span", "data-name")
		if err != nil {
			return nil, errors.New("Error to get id server: " + err.Error())
		}

		srv.ID, err = strconv.Atoi(string(id[0]))
		if err != nil {
			return nil, err
		}

		srv.Name = string(htmlutils.GetNodeText(server, "span"))

		s = append(s, srv)
	}

	return s, nil
}