<script>
  export let value = {};
  export let label = "";
</script>

<div class="textarea-container">
  <div class="textarea-label">{label}</div>
  <textarea class="textarea-box input-area" type="" bind:value />
</div>