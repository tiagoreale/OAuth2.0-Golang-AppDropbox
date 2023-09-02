# OAuth2.0-Golang-AppDropbox

Hoje, vou abordar a integração de mecanismos de acesso OAuth2 para a API do Dropbox, utilizando o conceito de 'refresh token' para obter um token de acesso temporário. Para ilustrar, demonstrarei como exibir um token na tela, obtido a partir de um arquivo temporário de armazenamento de chaves gerado após a execução da aplicação 'app-upload.go', utilizando a linguagem de programação Go.
Esta integração é parte de um sistema de automação desenvolvido por Tiago Reale, que efetua backup e envio de arquivos para o Dropbox, com base em diversos cenários, seja em ambientes locais ou na nuvem. No nosso exemplo, estamos adotando a terminação '.ZIP' para os arquivos, mas isso pode ser facilmente modificado para se adequar a várias extensões de arquivo.
Durante o processo, registraremos informações de progresso e enviaremos notificações para o Telegram, além de armazenar detalhes em um arquivo de log.
