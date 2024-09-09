<script lang="ts">
  export let value: number = 0;
  export let label = "";
  export let max: number = undefined;
  export let min: number = undefined;
  let bindVal: number = 0;

  function handleInput(e: any) {
    let targetVal = Number(e.target.value);
    if (max !== undefined && targetVal >= max) {
      value = max;
      bindVal = value;
      return;
    }
    if (max !== min && targetVal <= min) {
      value = min;
      bindVal = value;
      return;
    }
    value = targetVal;
  }
</script>

<div class="input_container">
  <div class="input_label">{label}</div>
  <input
    class="input_box input_area"
    type="number"
    on:change={handleInput}
    bind:value={bindVal}
    {min}
    {max}
  />
</div>
