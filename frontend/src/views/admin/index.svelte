<script>
  import { navigateTo } from "svelte-router-spa";

  let buttons = [
    {
      url: "/admin/dishes",
      text: "блюда",
    },
    {
      url: "/dishes",
      text: "заказать",
    },
  ];
</script>

<div class="grid grid-cols-1 medium-grid-cols-3 grid-gap">
  {#each buttons as item}
    <button class="btn" on:click={() => navigateTo(item.url, "", true)}
      >{item.text}</button
    >
  {/each}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }
  .btn {
    width: 40vw;
    height: 10vh;
  }
</style>
