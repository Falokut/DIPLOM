<script>
  export let value = {};
  export let label = "";
</script>

<div class="input_container">
  <div class="input_label">{label}</div>
  <input class="input_box input_area" type="text" bind:value />
</div>
