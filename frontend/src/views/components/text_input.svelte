<script>
  export let value = {};
  export let label = "";
  export let onChange = ()=>{};
</script>

<div class="input-div">
  {#if label.length > 0}
    <div class="input-label">{label}</div>
  {/if}
  <input class="input-box input-area" type="text" bind:value on:change={onChange}/>
</div>
