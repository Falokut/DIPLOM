<script>
  export let selected = [];
  export let options = [];
  export let onchange;
  export let label = "";
  import MultiSelect from "svelte-multiselect";
</script>

<div class="input_container">
  <div class="input_label">{label}</div>
  <div class="input_box">
    <MultiSelect
      bind:selected
      bind:options
      on:change={onchange}
      noMatchingOptionsMsg={"ничего не найдено"}
    />
  </div>
</div>

<style>
  :global(div.multiselect > ul.selected > li) {
    width: auto;
    overflow: hidden;
    max-height: 100%;
  }
  :global(div.multiselect > ul.selected) {
    max-height: 4em;
    overflow-y: scroll;
  }
</style>
