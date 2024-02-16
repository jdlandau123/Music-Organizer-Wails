<script>
  import logo from './assets/images/logo-universal.png';
  import Header from './components/Header.svelte';
  import { SyncMusicCollection } from '../wailsjs/go/main/App.js';

  function test() {
    SyncMusicCollection().then(err => console.log(err))
  }
</script>

<main>
  <Header />
  <button on:click={test}>Test</button>
</main>
