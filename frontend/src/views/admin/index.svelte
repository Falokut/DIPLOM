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
      url: "/admin/restaurants",
      text: "рестораны",
    },
    {
      url: "/orders",
      text: "Мои заказы",
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
