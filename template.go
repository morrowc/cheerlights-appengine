package main

var (
	reportTemplate = `
<html>
  <body>
    <center>
      <img src="http://static.tumblr.com/vhktvf6/aWQlvif7m/cheerlights_logo.png"></br>
      <script type="text/javascript">document.write(unescape("%3Cscript src='" + "http://www.iobridge.com/widgets/io.js?0SKzLy5H02MP' type='text/javascript'%3E%3C/script%3E"));</script>
    <br>Update Reporting
    </center>
    <hr>
    <center>
    <table>
      <tr>
        <th align=left>Color</th>
        <th>SrcIP</th>
        <th align=left>Time/Date</th>
      </tr>
    {{ range . }}
        <tr>
          <td width=30%>{{ .Color }}</td>
          <td width=20% align=center>{{ .SourceIp }}</td>
          <td width=30%>{{ .Date }}</td>
        </tr>
    {{ end }}
    </table>
    <hr>
  </center>
  </body>
</html>
`
)
