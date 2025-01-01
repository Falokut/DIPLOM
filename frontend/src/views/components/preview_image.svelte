<script>
  export let url = "";
  export let alt = "";
</script>

<img
  class="preview"
  src={url}
  {alt}
  style="width: 100%;height: 100%;"
/>

<style>
  .preview {
    object-fit: scale-down fill;
  }
</style>
