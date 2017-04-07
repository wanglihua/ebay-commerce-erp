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
	<div style="color:#929293;text-align:center;width:200px;margin-left:auto;margin-right:auto;margin-top:20px;margin-bottom:30px;"><a style="color:#929293;" href="http://www.miibeian.gov.cn/">浙ICP备17002068号-2	卧龙海外仓储</a></div>
</body>
</html>
`)
	return footer_buffer.String()
}

func GetSiteFooter() string {
	var footer_buffer bytes.Buffer

	footer_buffer.WriteString(`
	<div style="color:#929293;text-align:center;width:200px;margin-left:auto;margin-right:auto;margin-top:20px;margin-bottom:30px;"><a style="color:#929293;" href="http://www.miibeian.gov.cn/">浙ICP备17002068号-2	卧龙海外仓储</a></div>
</body>
</html>
`)
	return footer_buffer.String()
}
