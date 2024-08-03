<script lang="ts">
  import { onMount } from "svelte";
  import { GetCart, ProcessOrder } from "../../client/cart";
  import { GetDishes } from "../../client/dish";
  import CartItem from "./cart_item.svelte";
  import {
    initMainButton,
    retrieveLaunchParams,
  } from "@telegram-apps/sdk";
  const { initData } = retrieveLaunchParams();
  import { navigateTo } from "svelte-router-spa";

  const mainButtonRes = initMainButton();
  let mainButton = mainButtonRes[0];

  let dishes = [];
  onMount(async () => {
    let cartItems = GetCart();
    let dishesIds = [];
    cartItems.forEach((v, k) => {
      if (v == 0) return;
      dishesIds.push(k);
    });

    dishes = await GetDishes(dishesIds);
    
    mainButton.setParams({
      text: "заказать",
      isVisible: true,
    });
    mainButton.enable();
    mainButton.on("click", async () => {
      mainButton.disable();
      let result = await ProcessOrder(initData.user.id);
      mainButton.enable();
      if (!result) {
        return;
      }
      console.log("вы заказали еду");
      navigateTo("/");
    });
  });
</script>

<div>Корзина:</div>
<div>
  {#each dishes as dish}
    <CartItem bind:dishInst={dish} />
  {/each}
</div>
