<html>
  <body style="margin: 0 auto;max-width: 800px;font-family: Arial, Sans-serif;font-size: 15px;color: #555;">
    <h1 style="text-align: center;text-transform: uppercase;font-weight: normal;color: #333;margin-bottom: 30px;padding-bottom: 20px;border-bottom: 1px solid #EEE;">New form submission</h1>

    <p>Hi,</p>

    <p>Someone just submitted your form. Here's what they had to say:</p>
    <table style="width: 100%;border: 0;border-collapse: collapse;">
    {{ range $key, $element := . }}
      <tr>
        <td style="border-top: 1px solid #CCC;padding: 8px;color: #000;"><b>{{ $key }}</b></td>
        <td style="border-top: 1px solid #CCC;padding: 8px;color: #000;">{{ range $element }}{{ . }}{{ end }}</td>
      </tr>
    {{ end }}
    </table>

    <p>Thanks for using Crema Forms.</p>
  </body>
</html>
