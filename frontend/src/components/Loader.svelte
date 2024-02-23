<div class="wrapper">
  <div class="loader"></div>
</div>

<style>
  .wrapper {
    display: flex;
    justify-content: center;
    margin-top: 30px;
    width: 60%;
  }

  .loader {
    height: 100px;
    width: 100px;
    border: 10px solid #ccc;
    border-right-color: #4ba346;
    border-radius: 50%;
    animation: rotate 1.5s linear infinite;
  }

  @keyframes rotate{
    100% {
      transform: rotate(360deg);
    }
  }
</style>
