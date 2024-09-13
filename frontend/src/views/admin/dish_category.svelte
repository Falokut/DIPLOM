<script lang="ts">
  import { GetUserIdByTelegramId } from "../../client/user";
  import {
    RenameDishesCategory,
    DeleteDishesCategory,
  } from "../../client/dishes_categories";
  import { retrieveLaunchParams } from "@telegram-apps/sdk";
  const { initData } = retrieveLaunchParams();

  export let category = {
    name: "",
    id: 0,
  };
  export let remove = function () {};

  async function updateDishName() {
    let userId = await GetUserIdByTelegramId(initData.user.id);
    if (userId.length == 0) {
      return;
    }

    RenameDishesCategory(userId, category.name, category.id);
  }

  async function deleteCategory() {
    let userId = await GetUserIdByTelegramId(initData.user.id);
    if (userId.length == 0) {
      return;
    }

    await DeleteDishesCategory(userId, category.id);
    remove();
  }
</script>

<div class="category_container">
  <input
    class="input_area category_input"
    type="text"
    bind:value={category.name}
  />
  <div class="buttons_container">
    <button class="apply_button" on:click={() => updateDishName()}>✓</button>
    <button class="remove_button" on:click={() => deleteCategory()}>✕</button>
  </div>
</div>

<style>
  .category_container {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    font-size: large;
  }
  .category_input {
    border-radius: 3px;
  }
  .buttons_container {
    padding-top: 10px;
    padding-bottom: 10px;
    font-size: small;
    text-align: center;
    vertical-align: middle;
    width: auto;
    display: flex;
    flex-direction: row;
    font-weight: 100;
    width: 30vw;
    align-items: center;
    justify-content: space-around;
  }

  .remove_button {
    width: 40%;
    display: block;
    background-color: transparent;
    border: 1px solid var(--tg-theme-bg-color);
  }

  .remove_button:hover {
    color: var(--tg-theme-destructive-text-color);
  }

  .apply_button {
    width: 40%;
    display: block;
    background-color: transparent;
    border: 1px solid var(--tg-theme-bg-color);
  }
  .apply_button:hover {
    color: green;
  }
</style>
