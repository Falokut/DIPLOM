<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import DishPreview from "./dish_preview.svelte";
  import { GetDishesCategories } from "../../client/dishes_categories";
  import { AddDish } from "../../client/dish";
  import { retrieveLaunchParams } from "@telegram-apps/sdk";
  import { ToBase64 } from "../../utils/base64";
  import { navigate } from "svelte-routing";

  import { initMainButton, initBackButton } from "@telegram-apps/sdk";
  import TextInput from "../components/text_input.svelte";
  import NumInput from "../components/num_input.svelte";
  import ImageInput from "../components/image_input.svelte";
  import MultiSelectInput from "../components/multi_select_input.svelte";
  import { GetUserIdByTelegramId } from "../../client/user";
  import TextAreaInput from "../components/text_area_input.svelte";
  const defaultDish = {
    name: "",
    categories: [],
    description: "",
    url: "",
    price: 0,
  };
  let dish = defaultDish;
  let selectedCategories = [];
  let image = {};

  let dishesCategoriesMap: Map<string, number> = new Map<string, number>();
  const mainButtonRes = initMainButton();
  let mainButton = mainButtonRes[0];

  const backButtonRes = initBackButton();
  let backButton = backButtonRes[0];

  let removeMainButtonListFn;
  let removeBackButtonListFn;
  onMount(() => {
    mainButton.setParams({
      text: "Добавить блюдо",
      isVisible: true,
    });

    mainButton.enable();
    removeMainButtonListFn = mainButton.on("click", async () => {
      mainButton.disable();
      const { initData } = retrieveLaunchParams();
      if (initData === undefined || initData.user === undefined) {
        return false;
      }
      let userId = await GetUserIdByTelegramId(initData.user.id);
      if (userId.length == 0) {
        return false;
      }

      let imageBody = {};
      await ToBase64(image).then((data) => (imageBody = data));
      let req = {
        name: dish.name,
        description: dish.description,
        categories: dish.categories,
        price: dish.price,
        image: imageBody,
      };
      let ok = AddDish(req, userId);
      if (ok) {
        dish = defaultDish;
        selectedCategories = [];
      }
      mainButton.enable();
    });

    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/admin", { replace: true });
    });
  });
  onDestroy(() => {
    mainButton.hide();
    removeMainButtonListFn();
    removeBackButtonListFn();
  });

  async function loadDishesCategories() {
    let categories = await GetDishesCategories();
    let dishCategories = [];
    categories.forEach((value) => {
      dishCategories.push(value.name);
      dishesCategoriesMap.set(value.name, value.id);
    });
    return dishCategories;
  }
</script>

<h3>Добавить блюдо</h3>
<div class="input_container">
  <div class="input_div">
    <ImageInput
      bind:outputUrl={dish.url}
      label={"картинка:"}
      bind:file={image}
      uploadLabel={"выбрать файл"}
    />
  </div>
  <div class="input_div">
    <TextInput bind:value={dish.name} label={"название:"} />
  </div>
  <div class="input_div">
    <NumInput bind:value={dish.price} label={"цена:"} min={0} max={1000000} />
  </div>
  {#await loadDishesCategories() then dishCategories}
    <MultiSelectInput
      options={dishCategories}
      label={"категории:"}
      bind:selected={selectedCategories}
      onchange={() => {
        dish.categories = [];
        selectedCategories.forEach((name) => {
          let id = dishesCategoriesMap.get(name);
          dish.categories.push(id);
        });
      }}
    />
  {/await}

  <div>
    <TextAreaInput bind:value={dish.description} label={"описание"} />
  </div>
</div>

<h3>Предосмотр:</h3>
{#if dish.url != "" && dish.name != "" && dish.price > 0}
  <DishPreview bind:dish />
{/if}

<style>
  .input_container {
    display: flex;
    flex-flow: column;
    vertical-align: middle;
    justify-items: center;
    background-color: var(--tg-theme-secondary-bg-color);
    border-radius: 5px;
    height: 40vh;
  }
</style>
