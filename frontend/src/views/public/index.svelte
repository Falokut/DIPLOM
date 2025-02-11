<script>
  import { navigate } from "svelte-routing";

  let buttons = [
    {
      url: "/orders",
      text: "Мои заказы",
    },
    {
      url: "/dishes",
      text: "Заказать",
    },
  ];
</script>

<main class="btn-grid">
  {#each buttons as item}
    <button class="btn" on:click={() => navigate(item.url, { replace: true })}
      >{item.text}
    </button>
  {/each}
</main>

<style>
  .btn {
    width: 40vw;
    height: 20vh;
  }
</style>
