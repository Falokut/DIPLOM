<script lang="ts">
  import { DynamicScroll } from "svelte-dynamic-scroll";
  import { retrieveLaunchParams } from "@telegram-apps/sdk";
  import { GetUserOrders, UserOrder } from "../../client/order";
  import OrderItem from "./order_item.svelte";
  import { GetUserIdByTelegramId } from "../../client/user";
  const { initData } = retrieveLaunchParams();

  const pageLimit = 1;
  let currentOffset = 0;

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
</script>

<div class="orders_container">
  <DynamicScroll nextChunk={getMyOrders} let:value>
    <OrderItem order={value} />
  </DynamicScroll>
</div>

<style>
  .orders_container {
    height: calc(100vh - 50px);
    border-radius: 8px;
    width: 95vw;
    background-color: var(--tg-theme-secondary-bg-color);
  }
</style>
