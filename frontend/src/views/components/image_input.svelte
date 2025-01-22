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
    document.getElementById("file-upload-input").click();
  }
</script>

<div class="input-div">
  <div class="input-label">{label}</div>
  <div class="input-box">
    <button class="file-upload-button btn" on:click={openFileInput}
      >{uploadLabel}
    </button>
  </div>
  <input
    type="file"
    id="file-upload-input"
    style="display:none;"
    bind:files
    on:change={fileInput}
    accept="image/png, image/jpeg"
  />
</div>

<style>
  .file-upload-button {
    width: 100%;
  }
</style>
