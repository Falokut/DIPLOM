<script>
  export let file = {};
  export let label = "";
  export let outputUrl = "";
  export let uploadLabel = "";
  let files = [];
  function fileInput() {
    if (files.length != 1 || !files[0].type.match("image/*")) {
      files = file === undefined ? [] : [file];
      alert("невалидный файл с изображением");
      return;
    }
    file = files[0];
    outputUrl = URL.createObjectURL(file);
  }

  function openFileInput() {
    document.getElementById("file_upload_input").click();
  }
</script>

<div class="input_container">
  <div class="input_label">{label}</div>
  <button class="file_upload_button input_box" on:click={openFileInput}
    >{uploadLabel}</button
  >
  <input
    type="file"
    id="file_upload_input"
    style="display:none;"
    bind:files
    on:change={fileInput}
    accept="image/png, image/jpeg"
  />
</div>
