<script lang="ts">
  import { DynamicScroll } from "svelte-dynamic-scroll";
  import { GetUserOrders, UserOrder } from "../../client/order";
  import OrderItem from "./order_item.svelte";
  import { GetUserIdByTelegramId } from "../../client/user";

  import { initBackButton, retrieveLaunchParams } from "@telegram-apps/sdk";
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";

  const { initData } = retrieveLaunchParams();
  const pageLimit = 30;
  let currentOffset = 0;
  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];

  async function getMyOrders(
    lastVal: UserOrder | undefined
  ): Promise<UserOrder[]> {
    let userId = await GetUserIdByTelegramId(initData.user.id);
    if (userId.length == 0) {
      return;
    }
    orders = await GetUserOrders(userId, currentOffset, pageLimit);
    currentOffset += orders.length;
    return orders;
  }
  let orders: UserOrder[] = [];

  let removeBackButtonListFn = () => {};
  onMount(() => {
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
  });
  onDestroy(() => {
    removeBackButtonListFn();
  });
</script>

<div class="orders_container">
  <DynamicScroll nextChunk={getMyOrders} let:value>
    <OrderItem order={value} />
  </DynamicScroll>
</div>

<style>
  .orders_container {
    height: 95vh;
    border-radius: 8px;
    width: 95vw;
    background-color: var(--tg-theme-secondary-bg-color);
  }
</style>
