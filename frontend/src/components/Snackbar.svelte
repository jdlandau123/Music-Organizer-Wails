<script>
  import { fade } from 'svelte/transition';

  let isOpen = false;
  let msg = '';
  
  export function open(message) {
    isOpen = true;
    msg = message;
    setTimeout(() => {
      isOpen = false;
      msg = '';
    }, 3000);
  }
</script>

{#if isOpen}
  <div class="wrapper" transition:fade={{delay: 0, duration: 200}}>
    <p>{msg}</p>
  </div>
{/if}

<style>
  .wrapper {
    position: absolute;
    z-index: 1;
    bottom: 30px;
    left: 0;
    right: 0;
    margin-left: 40%;
    margin-right: 40%;
    padding: 10px;
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
    border-radius: 8px;
  }
</style>
