<script>
  export let url = "";
  export let alt = "";
  export let size = 90;
</script>

<img
  class="preview"
  src={url}
  {alt}
  style="width: {size}px;height: {size}px;"
/>

<style>
  .preview {
    object-fit: scale-down fill;
  }
</style>
