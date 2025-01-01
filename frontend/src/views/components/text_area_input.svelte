<script>
  export let value = {};
  export let label = "";
</script>

<div class="textarea_container">
  <div class="textarea_label">{label}</div>
  <textarea class="textarea_box input-area" type="" bind:value />
</div>
