package AuxClass.ProdutoClass;

public class Produto{
    private int codigo;
    private String nome;
    private int quantidade;
    private String tipo;
    private float valor;

    public Produto(int codigo){
        this.codigo = codigo;
        this.nome = "Nao definido";
        this.tipo = "Nao definido";
        this.quantidade = 0;
        this.valor = 0;
    }

    public Produto(int codigo,String nome){
        this.codigo = codigo;
        this.nome = nome;
        this.tipo = "Nao definido";
        this.quantidade = 0;
        this.valor = 0;
    }

    public Produto(int codigo,String nome,int quantidade){
        this.codigo = codigo;
        this.nome = nome;
        this.quantidade = quantidade;
        this.tipo = "Nao definido";
        this.valor = 0;
    }

    public Produto(int codigo,String nome,int quantidade,String tipo, float valor){
        this.codigo = codigo;
        this.nome = nome;
        this.quantidade = quantidade;
        this.tipo = tipo;
        this.valor = valor;
    }

    public void vender(int quantidade){
        if (this.quantidade >= quantidade){
            this.quantidade -= quantidade;
            float nota = quantidade * this.valor;
            System.out.println("Venda efetuada \nValor da Venda R$" + nota);
        }else{
            System.out.println("Quantidade insuficiente do produto em estoque ");
        }
    }

    public void comprar(int quantidade, float value){
        this.quantidade += quantidade;
        this.valor = value;

        System.out.println("Compra do produto efetuada com sucesso\n Quantidade atual :"+this.quantidade+"\nValor atual :"+this.valor);
    }

    @Override
    public String toString(){
        return "__________INFO PRODUTO__________\nCodigo do produto :"+this.codigo+"\nNome do produto :"+this.nome+"\nQuantidade em estoque :"+this.quantidade+"\nTipo do produto :"+this.tipo+"\nValor unitario do produto :"+this.valor;
    }

    public void inserir(String nome, int quantidade, int valor){
        this.nome = nome;
        this.quantidade = quantidade;
        this.valor = valor;
        System.out.println("Dados cadastrados do produto " + this.nome + "alterados");
    }


    public boolean igual(Produto produto){
        if ((produto.nome.equals(this.nome)) && (produto.tipo.equals(this.tipo))){
            return true;
        }else{
            return false;
        }
    }



}