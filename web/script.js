// STEP 1: select element and register change event
const imagePreview = document.querySelector('[data-target="image-preview"]');
const spinner = document.querySelector('[data-target="spinner"]');
const fileUploader = document.querySelector('[data-target="file-uploader"]');
fileUploader.addEventListener("change", handleFileUpload);

async function handleFileUpload(e) {
  try {
    const file = e.target.files[0];
    setUploading(true);
    if (!file) return;

    //const beforeUploadCheck = await beforeUpload(file);
    //if (!beforeUploadCheck.isValid) throw beforeUploadCheck.errorMessages;

    //const arrayBuffer = await getArrayBuffer(file);
    //const response = await uploadFileAJAX(arrayBuffer);

    //var formdata = new FormData();
    //formdata.append("file", e.files[0], e.files[0].name);

    var requestOptions = {
      method: 'POST',
      body: file,
    };

    fetch("http://localhost:3000/api/v0/plumb/file", requestOptions)
        .then(response => response.text())
        .then(result => console.log(result))
        .catch(error => console.log('error', error));
    
    alert("File Uploaded Success");
    //showPreviewImage(file);
  } catch (error) {
    alert(error);
    console.log("Catch Error: ", error);
  } finally {
    e.target.value = '';  // reset input file
    setUploading(false);
  }
}

// STEP 2: showPreviewImage with createObjectURL
// If you prefer Base64 image, use "FileReader.readAsDataURL"
function showPreviewImage(fileObj) {
  const image = URL.createObjectURL(fileObj);
  imagePreview.src = image;
}

// STEP 3: change file object into ArrayBuffer
function getArrayBuffer(fileObj) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    // Get ArrayBuffer when FileReader on load
    reader.addEventListener("load", () => {
      resolve(reader.result);
    });

    // Get Error when FileReader on error
    reader.addEventListener("error", () => {
      reject("error occurred in getArrayBuffer");
    });

    // read the blob object as ArrayBuffer
    // if you nedd Base64, use reader.readAsDataURL
    reader.readAsArrayBuffer(fileObj);
  });
}

// STEP 4: upload file throguth AJAX
// - use "new Uint8Array()"" to change ArrayBuffer into TypedArray
// - TypedArray is not a truely Array,
//   use "Array.from()" to change it into Array
function uploadFileAJAX(arrayBuffer) {
  // correct it to your own API endpoint
  return fetch("http://localhost:3000/api/v0/plumb/file", {
    headers: {
        "Content-Type": "application/octet-stream"
    },
    method: "POST",
    body: arrayBuffer
  })
    .then(res => {
      if (!res.ok) {
        throw res.statusText;
      }
      return res.json();
    })
    .then(data => data)
    .catch(err => console.log("err", err));
}

// STEP 5: Create before upload checker if needed
function beforeUpload(fileObject) {
  return new Promise(resolve => {
    const validFileTypes = ["image/jpeg", "image/png"];
    const isValidFileType = validFileTypes.includes(fileObject.type);
    let errorMessages = [];

    if (!isValidFileType) {
      errorMessages.push("You can only upload JPG or PNG file!");
    }

    resolve({
      isValid: isValidFileType,
      errorMessages: errorMessages.join("\n")
    });
  });
}

function setUploading(isUploading) {
  if (isUploading === true) {
    spinner.classList.add("opacity-1");
  } else {
    spinner.classList.remove("opacity-1");
  }
}