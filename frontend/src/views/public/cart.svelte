<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { GetCart, ProcessOrder } from "../../client/cart";
  import { GetDishes } from "../../client/dish";
  import CartItem from "./cart_item.svelte";
  import {
    initMainButton,
    initBackButton,
    retrieveLaunchParams,
  } from "@telegram-apps/sdk";
  const { initData } = retrieveLaunchParams();
  import { navigate } from "svelte-routing";

  const mainButtonRes = initMainButton();
  let mainButton = mainButtonRes[0];

  const backButtonRes = initBackButton();
  let backButton = backButtonRes[0];

  let dishes = [];
  let total = 0;
  let removeMainButtonListFn;
  let removeBackButtonListFn;
  let wishes = "";

  onMount(async () => {
    let cartItems = GetCart();
    let dishesIds = [];
    cartItems.forEach((v, k) => {
      if (v == undefined || k == undefined || v == 0) return;
      dishesIds.push(k);
    });

    dishes = await GetDishes(dishesIds);
    dishes.forEach((dish) => {
      let count = cartItems.get(dish.id.toString());
      if (count == 0) {
        return;
      }
      total += dish.price * count;
    });

    mainButton.setParams({
      text: "перейти к оплате",
      isVisible: true,
    });

    mainButton.enable();
    removeMainButtonListFn = mainButton.on("click", async () => {
      mainButton.disable();
      let result = await ProcessOrder(initData.user.id, wishes);
      mainButton.enable();
      if (!result) {
        return;
      }
      navigate("/", { replace: true });
    });

    backButton.show();
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
  });

  onDestroy(() => {
    removeMainButtonListFn();
    removeBackButtonListFn();
    mainButton.hide();
  });
</script>

<h3>Корзина:</h3>
<div class="line"></div>
<div class="cart">
  {#each dishes as dish}
    <CartItem bind:dishInst={dish} />
  {/each}
  <div class="line"></div>
</div>
<div class="total_block">
  <div class="total_text_l">Итого:</div>
  <div class="total_text_r">{(total / 100).toFixed(2) + "₽"}</div>
</div>
<div class="wishes">
  <textarea bind:value={wishes} class="wishes-input"></textarea>
</div>

<style>
  .cart {
    font-weight: normal;
    font-size: medium;
  }
  .wishes {
    width: 100vw;
    height: 30vh;
    text-align: center;
  }

  .wishes-input {
    background-color: var(--tg-theme-secondary-bg-color);
    color: var(--tg-theme-text-color);
    font-size: 14px;

    width: 80vw;
    height: 20vh;
    border: 0px;
    border-radius: 3%;
    resize: none;
    outline: none;
  }

  .line {
    border-radius: 30%;
    height: 5px;
    background-color: var(--tg-theme-secondary-bg-color);
  }

  .total_block {
    display: flex;
  }

  .total_text_l {
    text-align: left;
    padding-left: 10px;
    width: 50%;
  }
  .total_text_r {
    text-align: right;
    padding-right: 10px;
    width: 50%;
  }
</style>
