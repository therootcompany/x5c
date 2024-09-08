var FileViewer = {};
(function () {
    "use strict";

    FileViewer.init = function ({
        $fileInput,
        $dragzone,
        $textarea,
        $contents,
    }) {
        FileViewer._fileInput = $fileInput;
        FileViewer._dragzone = $dragzone;
        FileViewer._textarea = $textarea;
        FileViewer._contents = $contents;

        document.removeEventListener("dragover", FileViewer._$dragover);
        document.addEventListener("dragover", FileViewer._$dragover);

        document.removeEventListener("dragleave", FileViewer._$dragleave);
        document.addEventListener("dragleave", FileViewer._$dragleave);

        document.removeEventListener("drop", FileViewer._$drop);
        document.addEventListener("drop", FileViewer._$drop);
    };

    FileViewer.$readFirstFile = async function (event) {
        let $file = event.target.files[0];
        if (!$file) {
            return;
        }
        await FileViewer._decodeAndupdateContents($file);
    };

    FileViewer._$dragover = function (ev) {
        ev.preventDefault();
        if (FileViewer._dragzone) {
            FileViewer._dragzone.classList.add("dragover");
        }
    };

    FileViewer._$dragleave = function (ev) {
        if (FileViewer._dragzone) {
            FileViewer._dragzone.classList.remove("dragover");
        }
    };

    FileViewer._$drop = async function (ev) {
        ev.preventDefault();
        if (FileViewer._dragzone) {
            FileViewer._dragzone.classList.remove("dragover");
        }

        let $file = ev.dataTransfer.files[0];
        if (!$file) {
            return;
        }
        if (FileViewer._fileInput) {
            FileViewer._fileInput.files = ev.dataTransfer.files;
            FileViewer._fileInput.dispatchEvent(new Event("change"));
            return;
        }

        await FileViewer._decodeAndupdateContents($file);
    };

    FileViewer._decodeAndupdateContents = async function ($file) {
        let text = await FileViewer._decodeFileContents($file);
        if (FileViewer._contents) {
            FileViewer._contents.textContent = text;
        }
        if (FileViewer._textarea) {
            FileViewer._textarea.value = text;
            FileViewer._textarea.dispatchEvent(new Event("change"));
        }
    };

    FileViewer._decodeFileContents = async function (file) {
        function readFile(resolve) {
            let reader = new FileReader();
            reader.onload = function (e) {
                let buffer = e.target.result;
                resolve(buffer);
            };
            reader.onerror = function (e) {
                console.warn("[FileViewer] Error: failed to read file:");
                console.warn(e.message);
                let bytes = new Uint8Array(0);
                resolve(bytes.buffer);
            };
            reader.readAsArrayBuffer(file);
        }

        let buffer = await new Promise(readFile);
        let text;
        try {
            console.log(new Uint8Array(buffer));
            text = FileViewer._strictTextDecoder.decode(buffer);
        } catch (e) {
            text = FileViewer._bufferToHex(buffer);
        }
        return text;
    };

    FileViewer._strictTextDecoder = new TextDecoder("utf-8", {
        fatal: true,
    });

    FileViewer._bufferToHex = function (buffer) {
        let bytes = new Uint8Array(buffer);
        let hs = [];
        for (let b of bytes) {
            let h = b.toString(16);
            h = h.padStart(2, "0");
            hs.push(h);
        }
        let hex = hs.join("");
        return hex;
    };
})();
