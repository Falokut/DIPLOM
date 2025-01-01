<script>
  export let value = {};
  export let label = "";
</script>

<div class="input-div">
  {#if label.length > 0}
    <div class="input-label">{label}</div>
  {/if}
  <input class="input-box input-area" type="text" bind:value />
</div>
