<script>
  import logo from './assets/images/logo-universal.png'
  import { SelectDirectory, SetConfig, GetConfig } from '../wailsjs/go/main/App.js'
  import { onMount } from 'svelte';

  let resultText = "Please select a directory"
  let collectionPath;
  let devicePath;

  let config = {
    CollectionPath: '',
    DevicePath: ''
  };

  function setCollectionPath() {
    SelectDirectory().then(path => collectionPath = path);
  }

  function setDevicePath() {
    SelectDirectory().then(path => devicePath = path);
  }

  function setConfig() {
    SetConfig(collectionPath, devicePath).then(() => getConfig());
  }

  function getConfig() {
    GetConfig().then(c => config = c);
  }

  onMount(() => getConfig());
</script>

<main>
  <div class="result">Collection: {config.CollectionPath}</div>
  <div class="result">Device: {config.DevicePath}</div>
  <button class="btn" on:click={setCollectionPath}>Select Collection</button>
  <button class="btn" on:click={setDevicePath}>Select Device</button>
  <button class="btn" on:click={setConfig}>Set Config</button>
</main>

<style>
  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }
</style>
