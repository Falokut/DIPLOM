<script>
  import { GetUserIdByTelegramId } from "../../client/user";
  import { AddDishCategory } from "../../client/dishes_categories";
  import { retrieveLaunchParams } from "@telegram-apps/sdk";
  const { initData } = retrieveLaunchParams();

  export let categoryName = "";
  export let OnAdd = (category = { id: 0, name: "" }) => {};

  async function addDishCategory() {
    if (categoryName == "") {
      return;
    }
    let userId = await GetUserIdByTelegramId(initData.user.id);
    if (userId.length == 0) {
      return;
    }

    let dishCategory = await AddDishCategory(userId, categoryName);
    if (dishCategory.errorCode == undefined) {
      OnAdd(dishCategory);
      categoryName = "";
      return;
    }
    console.log(dishCategory);
  }
</script>

<div class="category_container">
  <input
    class="input_area category_input"
    type="text"
    bind:value={categoryName}
  />
  <button class="add_button" on:click={addDishCategory}>Добавить</button>
</div>

<style>
  .category_container {
    border-radius: 3%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    font-size: large;
    background-color: var(--tg-theme-secondary-bg-color);
    align-items: center;
    margin: auto;
    padding: 5px;
  }
  .category_input {
    border-radius: 3px;
  }

  .add_button {
    display: block;
    border: 1px solid var(--tg-theme-bg-color);
    margin-right: 5px;
  }
</style>
