package MainClass;
import AuxClass.ProdutoClass.Produto;

import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

public class Main{
    private int idlog = 1;
    List<Produto> produtos = new ArrayList<Produto>();


    public void main(String[] args){
        Scanner ler = new Scanner(System.in);
        int t = -1;
        Produto initTeste = new Produto(0,"Abertura de caixa",1,"Teste",0);
        while (t == -1){
            System.out.println("""
                                ____________________PDV PRODUTOS____________________
                                OPERAÇÕES :
                                1 = Cadastrar novo produto
                                2 = Visualizar dados de produto
                                3 = Atualizar dados de produto
                                4 = Comparar registro de produtos
                                5 = Registrar venda de produto   
                                6 = Registrar compra de produto
                                7 = Ver todos os produtos cadastrados
                                """);
            int option = ler.nextInt();

            switch (option) {
                case 1:
                    int l = -1;
                    while(l == -1)
                    System.out.println("""
                            Formas aceitaveis de cadastro :
                            1 = somente nome
                            2 = nome, quantidade
                            3 = nome, quantidade, tipo e valor  
                            """);
                    option = ler.nextInt();

                    switch (option) {
                        case 1:
                            System.out.println("Digite o nome a ser dado ao produto: ");
                            String nome = ler.nextLine();

                            int id = this.idlog;
                            Produto temp =  novoProduto(nome);
                            produtos.add(temp);

                            System.out.println("Produto cadastrado no id :" + id );
                            
                            break;

                        case 2:
                            int id2 = this.idlog;

                            System.out.println("Digite o nome a ser dado ao produto: ");
                            String nome2 = ler.nextLine();

                            System.out.println("Digite a quantidade do produto: ");
                            int qnt = ler.nextInt();

                            Produto temp2 =  novoProduto(nome2,qnt);
                            produtos.add(temp2);



                            System.out.println("Produto cadastrado no id :" + id2 );
                            
                            break;                        
                        case 3:
                            
                            int id3 = this.idlog;
                                
                            System.out.println("Digite o nome a ser dado ao produto: ");
                            String nome3 = ler.nextLine();

                            System.out.println("Digite a quantidade do produto: ");
                            int qnt2 = ler.nextInt();

                            System.out.println("Digite a categoria do produto: ");
                            String tipo = ler.nextLine();

                            System.out.println("Digite o valor do produto: ");
                            float vlt = ler.nextFloat();

                            Produto temp3 =  novoProduto(nome3,qnt2,tipo,vlt);
                            produtos.add(temp3);



                            System.out.println("Produto cadastrado no id :" + id3 );
                            break;
                    
                        default:
                            System.out.println("----------Opção inválida.----------");
                            break;
                    }
                    break;
            
                case 2:

                    break;

                case 3:
                    break;
                    
                case 4:
                    break;

                case 5:
                    break;

                case 6:
                    System.out.println();
                    break;

                default:
                    System.out.println("----------Opção invalida");
                    break;

            }
            
        }
        ler.close();


    }

    public Produto novoProduto(){
        Produto temp = new Produto(this.idlog);
        this.idlog += 1;
        return temp;
    }
    
    public Produto novoProduto(String nome){
        Produto temp = new Produto(this.idlog,nome);
        this.idlog += 1;
        return temp;
    }

    public Produto novoProduto(String nome, int quantidade){
        Produto temp = new Produto(this.idlog,nome,quantidade);
        this.idlog += 1;
        return temp;
    }
    
    public Produto novoProduto(String nome, int quantidade, String tipo, float valor){
        Produto temp = new Produto(this.idlog, nome, quantidade, tipo, valor);
        this.idlog += 1;
        return temp;
    }
    
}