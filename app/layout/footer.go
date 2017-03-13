package layout

import (
	"bytes"
)

func GetFooter() string {
	var footer_buffer bytes.Buffer

	footer_buffer.WriteString(`
        <script type="text/javascript">
            try {
                ace.settings.check('breadcrumbs', 'fixed')
            } catch (e) {
            }
        </script>
    </div>
</div>
</body>
</html>
`)
	return footer_buffer.String()
}

func GetSiteFooter() string {
	var footer_buffer bytes.Buffer

	footer_buffer.WriteString(`
</body>
</html>
`)
	return footer_buffer.String()
}
