package main

var html = []byte(`
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Upload</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/dropzone@5.4.0/dist/min/dropzone.min.css" media="none" onload="if(media!='all')media='all'">
    <style>
        body {
            background: #E8E9EC;
        }

        .dropzone {
            border: 2px dashed #0087F7;
            border-radius: 5px;
            background: white;
        }

        .dropzone .dz-message {
            font-weight: 400;
        }

        .dropzone .dz-message .note {
            font-size: 0.8em;
            font-weight: 200;
            display: block;
            margin-top: 1.4rem;
        }

        h1,
        h2,
        h3,
        table th,
        table th .header {
            font-size: 1.8rem;
            color: #0087F7;
            -webkit-font-smoothing: antialiased;
            line-height: 2.2rem;
        }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/dropzone@5.4.0/dist/min/dropzone.min.js"></script>
    <script>
        Dropzone.options.myAwesomeDropzone = {
            paramName: "file",
            maxFilesize: 999999,
            uploadMultiple: false,
            parallelUploads: 10,
            sending: function (file, xhr, formData) {
                formData.append("fullPath", file.fullPath || file.name);
            }
        };
    </script>
</head>

<body>
    <h1 id="try-it-out">Upload! <a href="/serve" target="_blank">Serve!</a></h1>
    <form action="/upload" class="dropzone needsclick dz-clickable" id="my-awesome-dropzone"></form>
</body>

</html>
`)
