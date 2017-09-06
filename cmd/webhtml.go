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
.building {
	margin-left: 20px;
}
.floor {
    margin-left: 40px;
}
.room {
    margin-left: 60px;
}
.row {
    margin-left: 80px;
}
.rack {
    margin-left: 100px;
}
        </style>
    </head>
    <body>
        <div id="nav">
<select>
    <option value="" style="display:none">Navigation</option>
    <option value="#apic">APIC</option>
    <option value="#fabricMembership">Fabric Membership</option>
    <option value="#geolocation">Geolocation</option>
</select>
        </div>
        <div id="apic">
            <a name="apic">
                <h2>APIC Authentication</h2>
            </a>
            <form method="POST">
                <label>URL:</label>
                <input class="apicInput" type="text" value="{{.URL}}" name="URL" readonly>
                <label>Username: </label>
                <input class="apicInput" type="text" value="{{.Username}}" name="Username" readonly="readonly">
                <label>Password: </label>
                <input class="apicInput" type="text" value="{{.Password}}" name="Password" readonly="readonly">
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
                {{- range $i, $e := .Nodes}}
                <div class="node">
                    <input class="fabricMembershipInput" type="text" value="{{$e.Name}}" name="Nodes.{{$i}}.Name" readonly>
                    <input class="fabricMembershipInput" type="text" value="{{$e.ID}}" name="Nodes.{{$i}}.ID" readonly>
                    <input class="fabricMembershipInput" type="text" value="{{$e.Serial}}" name="Nodes.{{$i}}.Serial" readonly>
                    <input class="fabricMembershipInput" type="text" value="{{$e.Pod}}" name="Nodes.{{$i}}.Pod" readonly>
                    <input class="fabricMembershipInput" type="text" value="{{$e.Role}}" name="Nodes.{{$i}}.Role" readonly> {{- end}}
                    <div class="buttons">
                        <input class="editButton" id="fabricMembershipEdit" name="Edit" type="button" value="Edit">
                        <input class="submitButton" type="submit" name="fabricMembershipSubmit" value="Submit">
                    </div>
                </div>
            </form>
        </div>
        <div id="geolocation">
            <a name="geolocation">
                <h2>Geolocation</h2>
            </a>
            <form method="POST">
		<p>{{.Sites}}</p>
                {{- range .Sites}}
                <div class="site">
                    <label>Site Name:</label>
                    <input class="geolocationInput" type="text" value="{{.Name}}" name="name" readonly>
                    <label>Description:</label>
                    <input class="geolocationInput" type="text" value="{{.Description}}" name="name" readonly>
                        {{- range .Buildings}}
                    <div class="building">
                        <label>Building Name:</label>
                        <input class="geolocationInput" type="text" value="{{.Name}}" name="name" readonly>
                        <label>Description:</label>
                        <input class="geolocationInput" type="text" value="{{.Description}}" name="name" readonly>
                        {{- range .Floors}}
                        <div class="floor">
                            <label>Floor Name:</label>
                            <input class="geolocationInput" type="text" value="{{.Name}}" name="name" readonly>
                            <label>Description:</label>
                            <input class="geolocationInput" type="text" value="{{.Description}}" name="name" readonly>
                        {{- range .Rooms}}
                            <div class="room">
                                <label>Room Name:</label>
                                <input class="geolocationInput" type="text" value="{{.Name}}" name="name" readonly>
                                <label>Description:</label>
                                <input class="geolocationInput" type="text" value="{{.Description}}" name="name" readonly>
                        {{- range .Rows}}
                                <div class="row">
                                    <label>Row Name:</label>
                                    <input class="geolocationInput" type="text" value="{{.Name}}" name="name" readonly>
                                    <label>Description:</label>
                                    <input class="geolocationInput" type="text" value="{{.Description}}" name="name" readonly>
                        {{- range .Racks}}
                                    <div class="rack">
                                        <label>Rack Name:</label>
                                        <input class="geolocationInput" type="text" value="{{.Name}}" name="rackName" readonly>
                                        <label>Description:</label>
                                        <input class="geolocationInput" type="text" value="{{.Description}}" name="rackDescription" readonly>
                                    </div>
                                    {{- end}}
                                </div>
                                {{- end}}
                            </div>
                            {{- end}}
                        </div>
                        {{- end}}
                    </div>
                    {{- end}}
                </div>
                {{- end}}
                <div class="buttons">
                    <input class="editButton" id="geolocationEdit" name="Edit" type="button" value="Edit">
                    <input class="submitButton" type="submit" name="geolocationSubmit" value="Submit">
                </div>
            </form>
        </div>
    </body>
    <script type="javascript"></script>
    <script>
var forms = [
    "apic",
    "fabricMembership",
    "geolocation"
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
