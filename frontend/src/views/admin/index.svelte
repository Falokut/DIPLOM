<script>
  import { navigate } from "svelte-routing";

  let buttons = [
    {
      url: "/admin/dishes/add",
      text: "Добавить блюдо",
    },
    {
      url: "/admin/dishes/delete",
      text: "Удалить блюдо",
    },
    {
      url: "/dishes",
      text: "Заказать",
    },
    {
      url: "/admin/dishes/categories",
      text: "Категории",
    },
    {
      url: "/orders",
      text: "Мои заказы",
    },
  ];
</script>

<main class="grid grid-cols-1 medium-grid-cols-3 grid-gap">
  {#each buttons as item}
    <button class="btn" on:click={() => navigate(item.url, { replace: true })}
      >{item.text}
    </button>
  {/each}
</main>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }
  .btn {
    width: 40vw;
    height: 20vh;
  }
</style>
