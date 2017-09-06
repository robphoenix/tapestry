package cmd

import "html/template"

var webTemplate = template.Must(template.New("web").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Tapestry</title>
<style>
body {
	font-family: sans-serif;
}
.buttons {
	display:block;
}
</style>
</head>
<body>
    <div id="nav">
        <ul>
            <li><a href="#apic">APIC</a></li>
            <li><a href="#fabricMembership">Fabric Membership</a></li>
            <li><a href="#geolocation">Geolocation</a></li>
        </ul>
    </div>
    <div id="apic">
        <a name="apic">
            <h2>APIC Authentication</h2>
        </a>
        <form method="POST">
            <label>URL:</label>
            <input class="apicInput" type="text" value="{{.URL}}" name="url" readonly>
            <label>Username: </label>
            <input class="apicInput" type="text" value="{{.Username}}" name="username" readonly="readonly">
            <label>Password: </label>
            <input class="apicInput" type="text" value="{{.Password}}" name="password" readonly="readonly">
            <div class="buttons">
                <input class="editButton" id="apicEdit" name="Edit" type="button" value="Edit">
                <input class="submitButton" type="submit" name="apicSubmit" value="Submit">
            </div>
        </form>
    </div>
    <div id="fabricMembership">
        <a name="fabricMembership">
            <h2>Fabric Membership</h2>
        </a>
        <form method="POST">
            {{- range .Nodes}}
            <div id="node{{.Name}}">
                <input class="fabricMembershipInput" type="text" value="{{.Name}}" name="name" readonly>
                <input class="fabricMembershipInput" type="text" value="{{.ID}}" name="id" readonly>
                <input class="fabricMembershipInput" type="text" value="{{.Serial}}" name="serial" readonly>
                <input class="fabricMembershipInput" type="text" value="{{.Pod}}" name="pod" readonly>
                <input class="fabricMembershipInput" type="text" value="{{.Role}}" name="role" readonly> {{- end}}
                <div class="buttons">
                    <input class="editButton" id="fabricMembershipEdit" name="Edit" type="button" value="Edit">
                    <input class="submitButton" type="submit" name="fabricMembershipSubmit" value="Submit">
                </div>
            </div>
        </form>
    </div>
</body>
<script type="javascript"></script>
<script>
var forms = [
    "apic",
    "fabricMembership"
];
forms.forEach(function(item, index, array) {
    console.log(item, index);
    document.getElementById(item + "Edit").onclick = function() {
        var inputs = document.getElementsByClassName(item + "Input");
        for (var i = 0; i < inputs.length; i++) {
            inputs[i].readOnly = false;
        }
    };
});
</script>
</html>
`))
