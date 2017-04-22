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
        <div class="beian"><a href="http://www.miibeian.gov.cn/">浙ICP备17002068号-2	卧龙海外仓储</a></div>
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
	<div class="beian"><a href="http://www.miibeian.gov.cn/">浙ICP备17002068号-2	卧龙海外仓储</a></div>
</body>
</html>
`)
	return footer_buffer.String()
}
